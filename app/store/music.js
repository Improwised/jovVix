import { defineStore } from "pinia";
export const useMusicStore = defineStore(
  "music-store",
  () => {
    const music = ref(false);

    const getMusic = () => {
      return music.value;
    };

    const setMusic = (data) => {
      music.value = data;
    };

    return { music, getMusic, setMusic };
  },
  {
    persist: true,
  }
);
