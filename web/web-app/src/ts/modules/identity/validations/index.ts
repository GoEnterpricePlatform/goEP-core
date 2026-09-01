import {
  isValidEmail,
  isValidPassword,
  isFormValidSignUp,
  isFormValidSignIn
} from "./validations";

// Optional: Export all for direct use
export {
  isValidEmail,
  isValidPassword,
  isFormValidSignUp,
  isFormValidSignIn
};

if (typeof window !== "undefined") {
  window.validationController = {
    isValidEmail,
    isValidPassword,
    isFormValidSignUp,
    isFormValidSignIn
  };
}
