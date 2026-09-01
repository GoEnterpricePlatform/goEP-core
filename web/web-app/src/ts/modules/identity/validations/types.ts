export interface ValidationController {
  isValidEmail: (email: string) => boolean;
  isValidPassword: (password: string) => boolean;
  isFormValidSignUp: (
    email: string,
    password: string,
    confirm_password: string,
  ) => boolean;
  isFormValidSignIn: (email: string, password: string) => boolean;
}

declare global {
  interface Window {
    validationController?: ValidationController;
  }
}
