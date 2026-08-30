import { helloWorld } from './hello';
import { helloName } from './hello-name';

// Optional: Export all for direct use
export { helloWorld, helloName };

// Exponer al window
if (typeof window !== "undefined") {
  window.greetingController = {
    helloWorld,
    helloName,
  };
}