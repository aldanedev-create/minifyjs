export function formatPrice(amount, currency) {
    const symbols = { USD: "$", EUR: "€", GBP: "£" };
    const symbol = symbols[currency] || currency;
    return symbol + amount.toFixed(2);
}

export function formatDate(date) {
    return date.toISOString().split("T")[0];
}