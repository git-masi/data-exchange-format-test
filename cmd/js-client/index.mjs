// @ts-check
import { exit } from "node:process";
import WebSocket from "ws";

(() => {
  const messageType = process.argv[2] || "bin";
  const host = process.argv[3] || "localhost:8080";

  switch (messageType) {
    case "bin":
      handleBin(`ws://${host}/bin`);
      break;
    case "json":
      handleJson(`ws://${host}/json`);
      break;
    default:
      console.log("No valid message type supplied");
      exit(1);
  }
})();

/**
 * @param {string} url
 */
function handleBin(url) {
  const ws = new WebSocket(url);
  let done = false;

  ws.on("close", () => {
    console.log("Test ended");
  });

  ws.on("error", (err) => {
    console.log("WebSocket error: ", err);
  });

  ws.on("open", () => {
    console.log("Test started");

    const buffer = new ArrayBuffer(16); // 8 bytes for BigInt64, 4 bytes each for two float32
    const view = new DataView(buffer);

    setTimeout(() => {
      done = true;
    }, 60_000);

    const interval = setInterval(() => {
      if (done) {
        clearTimeout(interval);
        ws.close(1000, "Test ended");
        return;
      }

      view.setBigInt64(0, BigInt(Date.now()), false);
      view.setFloat32(8, Math.random(), false);
      view.setFloat32(12, Math.random(), false);

      ws.send(buffer);
    }, 0);
  });
}

/**
 * @param {string} url
 */
function handleJson(url) {
  const ws = new WebSocket(url);
  let done = false;

  ws.on("close", () => {
    console.log("Test ended");
  });

  ws.on("error", (err) => {
    console.log("WebSocket error: ", err);
  });

  ws.on("open", () => {
    console.log("Test started");

    setTimeout(() => {
      done = true;
    }, 60_000);

    const interval = setInterval(() => {
      if (done) {
        clearTimeout(interval);
        ws.close(1000, "Test ended");
        return;
      }
      const message = {
        Time: Date.now(),
        Latitude: Math.random(),
        Longitude: Math.random(),
      };

      ws.send(JSON.stringify(message));
    }, 0);
  });
}
