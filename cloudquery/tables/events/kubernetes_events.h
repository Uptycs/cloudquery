#include "cloudquery/events/api_events_publisher.h"

namespace osquery {
namespace cloudquery {

class KubernetesEventSubscriber final
    : public EventSubscriber<APIEventPublisher> {
 public:
  Status init() override;

  Status KubernetesEventCallback(const ECRef& ec, const SCRef& sc);

  Status GetEventRows(std::vector<Row>& emitted_row_list,
                      std::string &sEventJson); 
};
}
} // namespace osquery
