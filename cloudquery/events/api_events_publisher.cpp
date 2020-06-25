#include "osquery/logger.h"
#include "cloudquery/events/api_events_publisher.h"
#include <osquery/registry_factory.h>
#include "osquery/remote/http_client.h"
#include "osquery/core/conversions.h"


namespace osquery {
REGISTER_CLOUDQUERY(APIEventPublisher, "event_publisher", "kubernetes_events");
namespace cloudquery {


void threadTrampoline(APIEventSubscriptionContextRef sub_ctx)
{
    osquery::http::Client client_; 
    osquery::http::Request request_(sub_ctx->sWatchURL);
    osquery::http::Response response_;

    auto tokens = osquery::split(sub_ctx->sWatchURL,"/");
    request_ << osquery::http::Request::Header("User-Agent", "osquery");
    request_ << osquery::http::Request::Header("Host", tokens[1]);
    for(auto header : sub_ctx->mHeaders) 
    {
        request_ << osquery::http::Request::Header(header.first, header.second);
        VLOG(1)<<header.first <<" = "<<header.second;
    }
    //pass on sub_ctx into lambda
    client_.get(request_, ([sub_ctx](std::string& sResponse) {
        auto ec = APIEventPublisher::createEventContext();
        ec->sub_ctx = sub_ctx;
        ec->sEventJson = sResponse;
        EventFactory::fire<APIEventPublisher>(ec);
    }));
    
}

void APIEventPublisher::configure()
{

    for (auto &sub : subscriptions_)
    {
        auto sub_ctx = getSubscriptionContext(sub->context);
        VLOG(1)<<" Subscribing to "<< sub_ctx->sWatchURL.c_str();
        std::shared_ptr<std::thread> threadObj(new std::thread(threadTrampoline, sub_ctx));
        VLOG(1)<<"started thread";
        sub_ctx->thread_ = threadObj;
    }
}

Status APIEventPublisher::run()
{

    pause(std::chrono::milliseconds(2000));
    return Status(0, "OK");
}

void APIEventPublisher::stop()
{

    for (auto &sub : subscriptions_)
    {
        auto sub_ctx = getSubscriptionContext(sub->context);
        sub_ctx->thread_.reset();
    }
}

bool APIEventPublisher::shouldFire(const APIEventSubscriptionContextRef& mc,
                                   const APIEventContextRef& ec) const
{
    return (mc == ec->sub_ctx);
    
}

}
}