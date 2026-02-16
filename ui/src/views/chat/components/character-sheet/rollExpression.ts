export type RollAdjustMode = 'normal' | 'adv' | 'dis';

const applyArgs = (template: string, args?: Record<string, any>): string => {
  let expr = template || '';
  for (const [key, value] of Object.entries(args || {})) {
    expr = expr.split('{' + key + '}').join(String(value));
  }
  return expr;
};

export const buildRollExpression = (
  template: string,
  args?: Record<string, any>,
  options?: { mode?: RollAdjustMode; modifier?: number },
): string => {
  const trimmed = applyArgs(template, args).trim();
  let command = '.ra';
  let body = trimmed;
  if (trimmed.startsWith('.')) {
    const firstSpace = trimmed.indexOf(' ');
    if (firstSpace === -1) {
      command = trimmed || '.ra';
      body = '';
    } else {
      command = trimmed.slice(0, firstSpace) || '.ra';
      body = trimmed.slice(firstSpace + 1);
    }
  }
  const bodyTrimmed = body.trim();
  const commandPrefix = command + ' ';
  if (bodyTrimmed.startsWith(commandPrefix)) {
    body = bodyTrimmed.slice(commandPrefix.length);
  } else if (bodyTrimmed === command) {
    body = '';
  }

  const modifierValue = options?.modifier ?? 0;
  const modText = modifierValue === 0 ? '' : (modifierValue > 0 ? '+' : '') + String(modifierValue);
  body = body.trim();
  const bodyWithMod = body && modText ? body + modText : body;
  const modeToken = options?.mode === 'adv' ? 'b' : options?.mode === 'dis' ? 'p' : '';

  if (!bodyWithMod) {
    const commandWithMod = command + modText;
    return modeToken ? (commandWithMod + ' ' + modeToken).trim() : commandWithMod.trim();
  }

  return modeToken ? (command + ' ' + modeToken + ' ' + bodyWithMod).trim() : (command + ' ' + bodyWithMod).trim();
};
