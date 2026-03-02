import js from "@eslint/js";
import globals from "globals";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import react from "eslint-plugin-react";
import jsxA11y from "eslint-plugin-jsx-a11y";
import unusedImports from "eslint-plugin-unused-imports";
import tseslint from "typescript-eslint";
import { defineConfig, globalIgnores } from "eslint/config";
export default defineConfig([
    globalIgnores(["dist"]),
    {
        files: ["**/*.{ts,tsx}"],
        extends: [
            js.configs.recommended,
            tseslint.configs.recommended,
            reactHooks.configs["recommended-latest"],
            reactRefresh.configs.vite,
        ],
        languageOptions: {
            ecmaVersion: 2020,
            globals: globals.browser,
            parserOptions: {
                ecmaFeatures: {
                    jsx: true,
                },
            },
        },
        plugins: {
            react,
            "jsx-a11y": jsxA11y,
            "unused-imports": unusedImports,
        },
        settings: {
            react: {
                version: "detect",
            },
        },
        rules: {
            "unused-imports/no-unused-imports": "error",
            "unused-imports/no-unused-vars": [
                "warn",
                {
                    vars: "all",
                    varsIgnorePattern: "^_",
                    args: "after-used",
                    argsIgnorePattern: "^_",
                },
            ],
            semi: ["error", "always"],
            quotes: ["error", "single", { avoidEscape: true }],
            "comma-dangle": ["error", "always-multiline"],
            "no-multiple-empty-lines": ["error", { max: 1, maxEOF: 0 }],
            "eol-last": ["error", "always"],
            "no-trailing-spaces": "error",
            "react/jsx-uses-react": "off",
            "react/react-in-jsx-scope": "off",
            "react/prop-types": "off",
            "react/jsx-curly-spacing": ["error", { when: "never", children: true }],
            "react/jsx-tag-spacing": [
                "error",
                {
                    closingSlash: "never",
                    beforeSelfClosing: "always",
                    afterOpening: "never",
                    beforeClosing: "never",
                },
            ],
            "@typescript-eslint/no-unused-vars": "off",
            "@typescript-eslint/explicit-module-boundary-types": "off",
            "@typescript-eslint/no-explicit-any": "warn",
            "jsx-a11y/alt-text": "warn",
            "jsx-a11y/anchor-is-valid": "warn",
        },
    },
]);
