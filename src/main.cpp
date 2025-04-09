#include <iostream>

#include "cli.h"
#include "command_runner.h"

int main(int argc, char** argv) {
  CLI cli("0.0.1");
  cli.run(argc, argv);
  CommandRunner cmdRunner;
  std::cout << cmdRunner.run("ls .") << std::endl;
  return 0;
}
