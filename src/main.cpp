#include "cli.h"

int main(int argc, char** argv) {
  CLI cli("0.0.1");
  cli.run(argc, argv);
  return 0;
}
