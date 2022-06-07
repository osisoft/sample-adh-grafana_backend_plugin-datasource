import { SelectableValue } from '@grafana/data';

export const debounce = (fn: (args: string) => Promise<Array<SelectableValue<string>>>, ms: number) => {
  let timer: NodeJS.Timeout;

  const debouncedFunc = (args: string): Promise<Array<SelectableValue<string>>> =>
    new Promise((resolve) => {
      if (timer) {
        clearTimeout(timer);
      }
      timer = setTimeout(() => {
        resolve(fn(args));
      }, ms);
    });

  return debouncedFunc;
};
