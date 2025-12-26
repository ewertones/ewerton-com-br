export function deepMerge(target: object, source: object): object {
  if (typeof target !== 'object' || target === null) {
    return source;
  }
  if (typeof source !== 'object' || source === null) {
    return target;
  }

  const output = { ...target } as Record<string, unknown>;

  for (const key of Object.keys(source)) {
    const sourceValue = (source as Record<string, unknown>)[key];
    const targetValue = (target as Record<string, unknown>)[key];
    
    if (sourceValue instanceof Object && key in target) {
      output[key] = deepMerge(targetValue as object, sourceValue as object);
    } else {
      output[key] = sourceValue;
    }
  }

  return output;
}
