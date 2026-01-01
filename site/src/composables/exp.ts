
export function getCurrentLevel(score: number) {
  let level = 0;
  while (getTotalExpForLevel(level + 1) <= (score || 0)) {
    level++;
  }
  return level + 1;
}

export function getLevelProgress(score: number) {
  const currentLevel = getCurrentLevel(score);
  const currentLevelExp = getTotalExpForLevel(currentLevel - 1);
  const nextLevelExp = getTotalExpForLevel(currentLevel);
  const expInCurrentLevel = score - currentLevelExp;
  const expNeededForNextLevel = nextLevelExp - currentLevelExp;

  return Math.min(100, Math.round((expInCurrentLevel / expNeededForNextLevel) * 100));
}

export function getExpToNextLevel(level: number) {
  const base = 10;
  const growth = 5;

  let exp = base + level * growth;

  if (level > 500) {
    exp = exp * Math.pow(1.01, level - 500);
  }

  return Math.floor(exp); // round down to integer
}

export function getTotalExpForLevel(level: number) {
  let total = 0;
  for (let i = 1; i <= level; i++) {
    total += getExpToNextLevel(i);
  }
  return total;
}