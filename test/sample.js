console.log("\nEnv vars:");
console.log("TEST_CLIENT_SECRET:", process.env.TEST_CLIENT_SECRET || "");
console.log("TEST_CLIENT_NAME:", process.env.TEST_CLIENT_NAME || "");
console.log("--------------------------------");

console.log("\nDoing some work...");
const barLength = 10;
let i = 1;
const interval = setInterval(() => {
  const filled = "█".repeat(i);
  const empty = "░".repeat(barLength - i);
  process.stdout.write(`\r[${filled}${empty}] ${i * 10}%`);
  if (i === barLength) {
    clearInterval(interval);
    console.log("\n✔ Done!\n");
  }
  i++;
}, 200);
