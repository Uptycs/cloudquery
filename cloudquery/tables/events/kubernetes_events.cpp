#include <osquery/logger.h>
#include <osquery/sql.h>

#include "cloudquery/tables/events/kubernetes_events.h"
#include <osquery/filesystem.h>
#include <osquery/registry_factory.h>
#include <sstream>
#include <iostream>

#define FILL_FROM_JSON(x,y,z) \
if(z.HasMember(#x) && z[#x].IsString())\
      y = z[#x].GetString();
namespace osquery {
  FLAG(string,
     kubernetes_api_server,
     "https://10.96.0.1/api/v1",
     "Kubernetes API server URL");

  FLAG(string,
     kubernetes_secret_path,
     "/run/secrets/kubernetes.io/serviceaccount/token",
     "Kubernetes service account token");

namespace cloudquery {

FLAG(bool,
     allow_kubernetes_events,
     true,
     "Allow Kubernetes Events");

Status KubernetesEventSubscriber::init() {
  if (!FLAGS_allow_kubernetes_events) {
    return Status(1, "Subscriber disabled via configuration");
  }

  auto sc = createSubscriptionContext();
  sc->sWatchURL = FLAGS_kubernetes_api_server + "/watch/events";
  std::string sSecret = "";
  Status status = readFile(FLAGS_kubernetes_secret_path, sSecret);
  if(status.ok() && !sSecret.empty())
  {
    sc->mHeaders["Authorization"] = "Bearer " + sSecret;
  }
  subscribe(&KubernetesEventSubscriber::KubernetesEventCallback, sc);
  return Status(0, "OK");
}

Status KubernetesEventSubscriber::GetEventRows(std::vector<Row> &emitted_row_list,
                    std::string &sEventJson)
{
  std::istringstream stream(sEventJson);
  std::string line;
  while (std::getline(stream, line) )
  {
    rapidjson::Document doc;
    Row r;
    if (doc.Parse(line.c_str()).HasParseError())
    {
      LOG(WARNING) << "Failed parsing row: " << line;
      continue;
    }
    FILL_FROM_JSON(type, r["type"], doc);
    if (doc.HasMember("object") && (doc["object"].IsObject()))
    {
      auto obj = doc["object"].GetObject();

      FILL_FROM_JSON(reason, r["reason"], obj);
      FILL_FROM_JSON(message, r["message"], obj);
      FILL_FROM_JSON(firstTimestamp, r["first_timestamp"], obj);
      FILL_FROM_JSON(lastTimestamp, r["last_timestamp"], obj);
      FILL_FROM_JSON(type, r["event_type"], obj);
      if (obj.HasMember("count") && obj["count"].IsInt())
      {
        r["count"] = BIGINT(obj["count"].GetInt());
      }
      if (obj.HasMember("metadata") && obj["metadata"].IsObject())
      {
        auto metaDataObj = obj["metadata"].GetObject();
        FILL_FROM_JSON(name, r["event_name"], metaDataObj);
        FILL_FROM_JSON(uid, r["event_uid"], metaDataObj);
        FILL_FROM_JSON(namespace, r["event_namespace"], metaDataObj);
        FILL_FROM_JSON(selfLink, r["self_link"], metaDataObj);
        FILL_FROM_JSON(creationTimestamp, r["timestamp"], metaDataObj);
      }
      if (obj.HasMember("involvedObject") && obj["involvedObject"].IsObject())
      {
        auto involvedObject = obj["involvedObject"].GetObject();
        FILL_FROM_JSON(kind, r["object_kind"], involvedObject);
        FILL_FROM_JSON(name, r["object_name"], involvedObject);
        FILL_FROM_JSON(uid, r["object_uid"], involvedObject);
      }
      if (obj.HasMember("source") && obj["source"].IsObject())
      {
        auto srcObj = obj["source"].GetObject();
        FILL_FROM_JSON(component, r["source_component"], srcObj);
        FILL_FROM_JSON(host, r["source_host"], srcObj);
      }
    }
    VLOG(1) << "EVENT got type = "<< r["type"] << " object_kind  "<<r["object_kind"];
    emitted_row_list.push_back(r);
  }
  return Status(0,"OK");
}

Status KubernetesEventSubscriber::KubernetesEventCallback(const ECRef& ec,
                                                          const SCRef& sc) {

  std::vector<Row> emitted_row_list;
  VLOG(1)<<ec->sEventJson;

  auto status = GetEventRows(emitted_row_list, ec->sEventJson);
  if (!status.ok()) {
    return status;
  }

   for (auto& row : emitted_row_list) {
     add(row);
   }

  return Status(0, "OK");
}

}
REGISTER_CLOUDQUERY(KubernetesEventSubscriber, "event_subscriber", "kubernetes_events");

} // namespace osquery
