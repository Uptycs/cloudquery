#include "osquery/logger.h"
#include "cloudquery/events/api_events_publisher.h"
#include <osquery/registry_factory.h>


namespace osquery {
REGISTER_CLOUDQUERY(APIEventPublisher, "event_publisher", "kubernetes_events");
namespace cloudquery {


Status APIEventPublisher::setUp()
{
    // Global curl init
    curl_global_init(CURL_GLOBAL_DEFAULT);
    return Status(0, "OK");

}

void threadTrampoline(CURL* curl)
{
    size_t start_time = getUnixTime();
    size_t end_time = 0;
    CURLcode curlCode = CURLcode::CURLE_OK;
    do{
        VLOG(1)<<"Started thread with curl " <<curl;
        curlCode = curl_easy_perform(curl);
        end_time = getUnixTime();

    } while (end_time - start_time > 2);

    VLOG(1)<<"Connection ended with curl handle  "<< curl << " error = "<< curlCode;
}

void APIEventPublisher::configure()
{

    for (auto &sub : subscriptions_)
    {
        auto sub_ctx = getSubscriptionContext(sub->context);
        CURL *curlPtr_ = curl_easy_init();
        if (!curlPtr_)
        {
            LOG(WARNING) << "Failed initializing curl";
            continue;
        }
        VLOG(1)<<" Subscribing to "<< sub_ctx->sWatchURL.c_str();
        SetCurlCallbackData(curlPtr_, sub_ctx);
        std::shared_ptr<std::thread> threadObj(new std::thread(threadTrampoline, curlPtr_));
        VLOG(1)<<"started thread";
        sub_ctx->curlPtr_ = curlPtr_;
        sub_ctx->thread_ = threadObj;
    }
}

void APIEventPublisher::SetCurlCallbackData(CURL *curl, APIEventSubscriptionContextRef sub_ctx)
{

    curl_easy_setopt(curl, CURLOPT_URL, sub_ctx->sWatchURL.c_str());
    struct curl_slist *list = NULL;

    for(auto header : sub_ctx->mHeaders) 
    {
        std::string headerString = header.first + ": " + header.second;
        list = curl_slist_append(list, headerString.c_str());
    }
    if(list) 
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, list);

    curl_easy_setopt(curl, CURLOPT_NOPROGRESS, 1L);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, APIEventPublisher::http_get_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void*)sub_ctx.get());
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
        curl_easy_cleanup(sub_ctx->curlPtr_);
    }
}

// callback that gets registered with curl
size_t APIEventPublisher::http_get_callback(char *responseChunk, size_t count, 
                                            size_t bytesPerCount, void *userdata)
{
    size_t bytesReceived = count * bytesPerCount;
    APIEventSubscriptionContext* sc = (APIEventSubscriptionContext*)userdata;

    std::string sResponse(responseChunk, bytesReceived);
    auto ec = createEventContext();
    ec->sub_ctx = sc;
    ec->sEventJson = sResponse;
    EventFactory::fire<APIEventPublisher>(ec);
    return bytesReceived;

}

bool APIEventPublisher::shouldFire(const APIEventSubscriptionContextRef& mc,
                                   const APIEventContextRef& ec) const
{
    return (mc.get() == ec->sub_ctx);
    
}

}
}