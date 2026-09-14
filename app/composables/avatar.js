const avatars = [
  "Sophia", "Jude", "Jade", "Ryan", "Adrian", "Chase", "Nolan",
  "Sadie", "Brian", "Aidan", "Destiny", "Maria", "Kingston", "Andrea",
  "Vivian", "Eden", "Wyatt", "Sawyer", "Jocelyn",
];

// Function to get a random avatar
export const getRandomAvatarName = () => {
  const randomIndex = Math.floor(Math.random() * avatars.length);
  return avatars[randomIndex];
};

// Function to get avatar URL by name
export const getAvatarUrlByName = (name) => {
  const seed = avatars.includes(name) ? name : "Eden";
  return `https://api.dicebear.com/10.x/clay/svg?seed=${encodeURIComponent(seed)}`;
};
