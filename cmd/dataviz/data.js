// @ts-check

(async () => {
  try {
    const text = await (await fetch("http://127.0.0.1:8080/metrics")).text();

    const lines = text.split("\n");

    const messageTypes = [
      "binary_message_latency_microseconds_bucket",
      "json_message_latency_milliseconds_bucket",
    ];

    const timeUnits = {
      0: "microseconds",
      1: "milliseconds",
    };

    const histogramData = Object.fromEntries(
      messageTypes
        .map((type) => lines.filter((str) => str.startsWith(type)))
        .map((l, i) => parseLines(l, timeUnits[i]))
        .map((dataset, i) => [messageTypes[i], dataset])
    );

    console.log(histogramData);
  } catch (error) {
    console.error(error);
  }
})();

/**
 * @param {Array<string>} lines
 * @param {string} unitOftime
 */
function parseLines(lines, unitOftime) {
  const data = /** @type {Array<Array<string | number>>}} */ ([
    [`Latency (${unitOftime})`, "Number of Requests"],
  ]);

  lines.forEach((l) => {
    const [description, numRequests] = l.split(" ");

    const time = description.match(/\d+|\+Inf/);

    if (!time) {
      throw new Error("Could not parse description");
    }

    data.push([time[0], +numRequests]);
  });

  return data;
}
