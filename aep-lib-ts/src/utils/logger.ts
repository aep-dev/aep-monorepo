import { Logger } from "tslog";

export const logger = new Logger({
  name: "aep-lib-ts",
  minLevel: 2, // 0: silly, 1: trace, 2: debug, 3: info, 4: warn, 5: error, 6: fatal
  type: "pretty",
});

export default logger;
