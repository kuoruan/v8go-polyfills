module.exports = {
  root: true,
  env: {
    es6: true,
    browser: false,
    node: true,
  },
  extends: ["eslint:recommended", "plugin:prettier/recommended"],
  plugins: ["prettier"],
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: "module",
  },
};
