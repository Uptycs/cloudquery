#include <string>
#include "osquery/core/conversions.h"
#include "osquery/extensions/interface.h"
#include "cloudquery/loadPlugins.h"
#include "cloudquery/version.h"
#include <osquery/filesystem.h>

#ifdef __linux__
#include <dlfcn.h>
#include <sys/select.h>
#include "osquery/core/flagalias.h"

#endif
namespace fs = boost::filesystem;

namespace osquery
{

CLI_FLAG(string,
         inproc_autoload,
         CLOUDQUERY_HOME "inproc_autoload",
         "Optional path to a list of autoloaded inproc extensions");

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

bool isExtSafe(std::string& path) {
  boost::trim(path);
  // Resolve acceptable extension binaries from autoload paths.
  if (isDirectory(path).ok()) {
    VLOG(1) << "Cannot autoload extension from directory: " << path;
    return false;
  }
  // Only autoload file which were safe at the time of discovery.
  // If the binary later becomes unsafe (permissions change) then it will fail
  // to reload if a reload is ever needed.
  fs::path extendable(path);
  // Set the output sanitized path.
  path = extendable.string();
  if (!pathExists(path).ok()) {
    LOG(WARNING) << "Extension doesn't exist at: " << path;
    return false;
  }
  if (!safePermissions(extendable.parent_path().string(), path, true)) {
    LOG(WARNING) << "Will not autoload extension" 
                 << " with unsafe directory permissions: " << path;
    return false;
  }

  VLOG(1) << "Found autoloadable extension " << path;
  return true;
}

bool loadInProcPlugins()
{
   std::string autoload_paths;
   if (!readFile(FLAGS_inproc_autoload, autoload_paths).ok())
   {
      VLOG(1)<<" Failed reading: " << FLAGS_inproc_autoload;
      return false;
   }

   // The set of binaries to auto-load, after safety is confirmed.
   std::set<std::string> autoload_binaries;
   for (auto &path : osquery::split(autoload_paths, "\n"))
   {
      if (isDirectory(path))
      {
         std::vector<std::string> paths;
         listFilesInDirectory(path, paths, true);
         for (auto &embedded_path : paths)
         {
            if (isExtSafe(embedded_path))
            {
               autoload_binaries.insert(std::move(embedded_path));
            }
         }
      }
      else if (isExtSafe(path))
      {
         autoload_binaries.insert(path);
      }
   }

   for (const auto &binary : autoload_binaries)
   {
      loadPluginByPath(binary);
   }
   return true;

}

} // namespace osquery



