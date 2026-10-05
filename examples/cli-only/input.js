// A small file used by the cli-only examples.
// Both minify.sh and optimize.sh read this.

function calculateTotal(price, tax) {
    const subtotal = price + (price * tax);
    const rounded = Math.round(subtotal * 100) / 100;
    return rounded;
}

function formatPrice(amount, currency) {
    const symbols = {
        USD: "$",
        EUR: "€",
        GBP: "£",
    };
    const symbol = symbols[currency] || currency;
    return symbol + amount.toFixed(2);
}

const total = calculateTotal(100, 0.15);
console.log(formatPrice(total, "USD"));