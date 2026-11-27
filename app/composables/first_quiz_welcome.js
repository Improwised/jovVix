// One-time "start your first quiz" hint after sign-up. The register form
// sets the flag before it posts to Kratos; the home page takes it once.
// sessionStorage keeps it to this tab and survives the Kratos redirects.
const KEY = "jv-new-signup";

export const markNewSignup = () => {
  try {
    sessionStorage.setItem(KEY, "1");
  } catch {
    // Storage blocked (private mode, site data off): skip the hint.
  }
};

export const takeNewSignup = () => {
  try {
    const isNew = sessionStorage.getItem(KEY) === "1";
    sessionStorage.removeItem(KEY);
    return isNew;
  } catch {
    return false;
  }
};
