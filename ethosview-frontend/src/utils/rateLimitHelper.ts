// Rate limit helper for API calls
let lastCallTime = 0;
const MIN_INTERVAL = 500; // Minimum 500ms between calls

export function withRateLimit<T extends any[], R>(
  fn: (...args: T) => Promise<R>
): (...args: T) => Promise<R> {
  return async (...args: T): Promise<R> => {
    const now = Date.now();
    const timeSinceLastCall = now - lastCallTime;
    
    if (timeSinceLastCall < MIN_INTERVAL) {
      const delay = MIN_INTERVAL - timeSinceLastCall;
      await new Promise(resolve => setTimeout(resolve, delay));
    }
    
    lastCallTime = Date.now();
    return fn(...args);
  };
}

export function createStaggeredCall<T>(
  calls: (() => Promise<T>)[],
  staggerMs = 200
): Promise<T[]> {
  return new Promise((resolve) => {
    const results: T[] = [];
    let completed = 0;
    
    calls.forEach((call, index) => {
      setTimeout(async () => {
        try {
          const result = await call();
          results[index] = result;
        } catch (error) {
          console.warn(`Staggered call ${index} failed:`, error);
          results[index] = null as any;
        }
        
        completed++;
        if (completed === calls.length) {
          resolve(results);
        }
      }, index * staggerMs);
    });
  });
}
