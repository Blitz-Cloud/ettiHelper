/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{html,js}"],
  theme: {
    extend: {},
  },
  plugins: [
    require("@catppuccin/tailwindcss")({
      prefix: "ctp",
      defaultFlavour: "latte",
    }),
  ],
};
