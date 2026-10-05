import { formatPrice } from "./utils.js";
import { add, multiply } from "./math.js";

const subtotal = add(10, 20);
const total = multiply(subtotal, 1.15);

console.log("subtotal:", formatPrice(subtotal, "USD"));
console.log("total:   ", formatPrice(total, "USD"));