# 📘 CSV File Instruction Set for Jovvix Quiz Platform

This guide provides rules and formatting instructions for uploading quiz questions using a CSV file in the Jovvix system.

---

## CSV Header Structure

Currently, your CSV file should include the following **12 headers**:

1. `Question Text`
2. `Question Type`
3. `Points`
4. `Option 1` to `Option 5`
5. `Correct Answer`
6. `Question Media`
7. `Option Media`
8. `Resource`

These headers **do not need to follow a strict order** — they can be rearranged as needed.

---

## Question Text Rules

- Do **not** use any quotes:
  -  `" "` or `' '` (double/single quotes)
- Do **not** use special characters:
  -  `:`, `-`, `,`, `\`, `/`
- Any language is supported (e.g., English, Hindi, Gujarati, etc.)
- Capital letters are allowed

---

## Question Type

There are **2 supported question types**:

1. **`single answer`**: Only one option is correct.
2. **`survey`**: All entered options are considered correct.

---

## Media-Based Questions

### Image Questions

- If the **question** contains an image:
  - Use the `Question Media` column to enter the image path or name.
- If the **options** include images:
  - Use the `Option Media` column.

**Important**: Only the following image formats are supported:
- `JPG`, `JPEG`, `PNG`, `WEBP`
- `SVG` is **not supported**

---

## Code-Based Questions

### Where to place code:

- Code in the question → `Question Media` column
- Code in the options → `Option Media` column
- Additional code, metadata, or explanation → `Resource` column

Use Markdown formatting or indentation to separate code clearly if needed.

---

## Points

- If `Points` column is left blank, it defaults to **1 point**.
- You may override by setting a custom point value.
- And our quiz platform allows maximum 20 points only.
---

## Options

- Maximum **5 options** (`Option 1` to `Option 5`)
- You can use fewer options (e.g., only 3) if desired
- Do not exceed 5 total options

---

## Correct Answer Formatting

### Single Answer Example:
```text
3
```
### Multiple Answers (Survey or multi-select):

```text
1|2|4
```

- Use the | (pipe symbol) to separate multiple correct option numbers.
