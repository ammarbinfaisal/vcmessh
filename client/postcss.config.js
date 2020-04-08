const tailwindcss = require("tailwindcss");
const autoprefixer = require("autoprefixer");
const purgecss = require("@fullhuman/postcss-purgecss");

const plugins = [tailwindcss, autoprefixer, ]

module.exports = { plugins };
