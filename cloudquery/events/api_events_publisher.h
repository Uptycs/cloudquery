/**
 * Copyright (c) 2020 Uptycs, Inc. All rights reserved
 */

#pragma once

#include "osquery/events.h"
#include "osquery/status.h"

#include <mutex>


namespace osquery
{
namespace cloudquery {

using HeaderMap = std::map<std::string, std::string>;
struct APIEventSubscriptionContext : public SubscriptionContext {
  std::string sWatchURL;
  HeaderMap mHeaders;
  std::shared_ptr<std::thread> thread_;

};
using APIEventSubscriptionContextRef = std::shared_ptr<APIEventSubscriptionContext>;

struct APIEventContext : public EventContext {
    APIEventSubscriptionContextRef sub_ctx;
    std::string sEventJson;
};
using APIEventContextRef = std::shared_ptr<APIEventContext>;


/**
 * @brief Event publisher for API Events.
 *
 * This EventPublisher allows EventSubscriber's to subscribe to keep-alive restful
 * API' streams using curl callback which keeps getting called for each new event.
 * For example it Watches 'kubernetes/event/watch' kind of URLs for k8s events
 */

class APIEventPublisher
    : public EventPublisher<APIEventSubscriptionContext, APIEventContext> {
  DECLARE_PUBLISHER("api_events");

 public:
  virtual ~APIEventPublisher() {
    stop();
  }
  // Filter to enum and notify Subscribers
  bool shouldFire(const APIEventSubscriptionContextRef& mc,
                  const APIEventContextRef& ec) const override;

  Status run() override;
  void configure() override;
  void stop() override;
  


 private:
};

}
} // namespace osquery
