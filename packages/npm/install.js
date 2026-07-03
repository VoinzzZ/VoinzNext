process.env.VOINZNEXT_INSTALL = "1";
require("./cli.js").downloadLatest().catch(() => {});