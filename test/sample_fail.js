console.log("--------------------------------");
console.log("\nStarting Node.js failure sample...");
console.log("Env vars:");
console.log("TEST_CLIENT_SECRET:", process.env.TEST_CLIENT_SECRET || "");
console.log("TEST_CLIENT_NAME:", process.env.TEST_CLIENT_NAME || "");

console.log("\nSimulating work before failure...");
const barLength = 5;
let i = 1;
const interval = setInterval(() => {
  const filled = "█".repeat(i);
  const empty = "░".repeat(barLength - i);
  process.stdout.write(`\r[${filled}${empty}] ${i * 20}%`);
  if (i === barLength) {
    clearInterval(interval);
    console.log("\n\n✘ Something went wrong — exiting with code 1");
    process.exit(1);
  }
  i++;
}, 200);
