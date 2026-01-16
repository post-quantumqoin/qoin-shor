#ifndef SOME_H
#define SOME_H

#include <string>
#include <utility>
#include <vector>

namespace Some {
	struct DeviceInfo {
	  std::string id;
	};

	std::vector<std::pair<std::string, std::string>> Generate(const DeviceInfo& device);
}

#endif  // SOME_H
