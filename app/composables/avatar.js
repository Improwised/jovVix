const avatars = [
  {
    name: "Sophia",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Sophia",
  },
  {
    name: "Jude",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Jude",
  },
  {
    name: "Jade",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Jade",
  },
  {
    name: "Ryan",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Ryan",
  },
  {
    name: "Adrian",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Adrian",
  },
  {
    name: "Chase",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Chase",
  },
  {
    name: "Nolan",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Nolan",
  },
  {
    name: "Sadie",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Sadie",
  },
  {
    name: "Brian",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Brian",
  },
  {
    name: "Aidan",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Aidan",
  },
  {
    name: "Destiny",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Destiny",
  },
  {
    name: "Maria",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Maria",
  },
  {
    name: "Kingston",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Kingston",
  },
  {
    name: "Andrea",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Andrea",
  },
  {
    name: "Vivian",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Vivian",
  },
  {
    name: "Eden",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Eden",
  },
  {
    name: "Wyatt",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Wyatt",
  },
  {
    name: "Sawyer",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Sawyer",
  },
  {
    name: "Jocelyn",
    url: "https://api.dicebear.com/9.x/bottts/svg?seed=Jocelyn",
  },
];

// Function to get a random avatar
export const getRandomAvatarName = () => {
  const randomIndex = Math.floor(Math.random() * avatars.length);
  return avatars[randomIndex].name;
};

// Function to get avatar URL by name
export const getAvatarUrlByName = (name) => {
  const avatar = avatars.find((avatar) => avatar.name === name);
  return avatar
    ? avatar.url
    : "https://api.dicebear.com/9.x/bottts/svg?seed=Eden";
};
