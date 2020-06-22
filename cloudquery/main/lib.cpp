#include <string>
#include "osquery/core/conversions.h"
#include "osquery/extensions/interface.h"
#include "cloudquery/loadPlugins.h"
#include "cloudquery/version.h"
#ifdef __linux__
#include <dlfcn.h>
#include <sys/select.h>

#endif

namespace osquery
{

namespace cloudquery
{


#ifdef DEBUG
const std::string kCQVersion = CONCAT(CLOUDQUERY_BUILD_VERSION, -debug);
#else
const std::string kCQVersion = STR(CLOUDQUERY_BUILD_VERSION);
#endif
const std::string kCQSDKVersion = CLOUDQUERY_SDK_VERSION;


} // namespace cloudquery

// Method to load inproc plugins from
// inprocPlugins.autoload
void loadPluginByPath(std::string sPluginPath)
{
   void *lib_handle;
   char *error;
   lib_handle = dlopen(sPluginPath.c_str(), RTLD_NOW | RTLD_GLOBAL);
   if (!lib_handle)
   {
      LOG(ERROR) << "Failed loading plugin " << sPluginPath.c_str() << dlerror();
      return;
   }
   LoadFn loadFn = (LoadFn)dlsym(lib_handle, "Load");
   if ((error = dlerror()) != NULL)
   {
      LOG(ERROR) << "Failed loading Load() function " << error;
      return;
   }

   ProvideRegistryBroadcastFn giveRegistryFn = (ProvideRegistryBroadcastFn)dlsym(lib_handle, "ProvideRegistryBroadcast");
   if ((error = dlerror()) != NULL)
   {
      LOG(ERROR) << "Failed loading ProvideRegistryBroadcast function " << error;
      return;
   }

   loadFn();
   RegistryBroadcast broadcast;
   giveRegistryFn(broadcast);
   RegistryFactory::get().addBroadcast((RouteUUID)lib_handle, broadcast, true);
}

bool loadInProcPlugins()
{
   loadPluginByPath(std::string("/home/batman/osquery/build/xenial/osquery/examples/in_proc_extensions/go_extensions/libexample_go_extension.so"));
   return true;
}

} // namespace osquery



