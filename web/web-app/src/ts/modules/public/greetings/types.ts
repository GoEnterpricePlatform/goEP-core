export interface GreetingController {
  helloWorld: () => void;
  helloName: (name: string) => void;
}

declare global {
  interface Window {
    greetingController?: GreetingController;
  }
}