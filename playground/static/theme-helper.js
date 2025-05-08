/**
 * Initialize the theme helper and set the theme stored
 * in the local storage
 *
 * It sets two global functions:
 * - getTheme: Returns the current theme
 * - setTheme: Sets the theme and stores it in the local storage
 */
function initThemeHelper() {
  function getTheme() {
    const theme = localStorage.getItem("theme");
    return theme || "";
  }

  function setTheme(theme) {
    const themeToSet = theme === "system" ? "" : theme;
    localStorage.setItem("theme", themeToSet);
    document.documentElement.setAttribute("data-theme", themeToSet);
  }

  globalThis.getTheme = getTheme;
  globalThis.setTheme = setTheme;

  const theme = getTheme();
  setTheme(theme);
}

initThemeHelper();
