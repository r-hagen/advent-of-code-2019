const fs = require('fs');

const fuel = function(m) {
  return Math.floor(m / 3) - 2;
}

const fuelR = function(m) {
  const f = fuel(m);
  return f < 0 ? 0 : f + fuelR(f);
}

const mm = fs.readFileSync('in').toString().split('\n').map(x => Number(x) || 0);

const ans1 = mm.reduce((f, m) => f + fuel(m), 0)
console.log(ans1);

const ans2 = mm.reduce((f, m) => f + fuelR(m), 0)
console.log(ans2);
