// Promise.withResolvers polyfill for older Node.js versions
if (typeof Promise.withResolvers !== 'function') {
  Promise.withResolvers = function<T>() {
    let resolve!: (value: T | PromiseLike<T>) => void;
    let reject!: (reason?: unknown) => void;
    
    // Create a new promise that exposes its resolver and rejecter
    const promise = new Promise<T>((res, rej) => {
      resolve = res;
      reject = rej;
    });
    
    return {
      promise,
      resolve,
      reject
    };
  };
}

export {}; // This makes the file a module 