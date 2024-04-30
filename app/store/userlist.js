import { defineStore } from "pinia";
export const useListUserstore = defineStore("listusers-store", () => {
  const listUsers = ref([]);

  //actions
  const addUser = (user) => {
    const checkUsers = listUsers.value.find(
      (currentUser) => currentUser === user
    );
    if (!checkUsers) {
      listUsers.value.unshift(user);
    }
  };
  const removeUser = (user) => {
    listUsers.value = listUsers.value.filter(
      (currentUser) => currentUser !== user
    );
  };
  const removeAllUsers = () => {
    listUsers.value = [];
  };

  return { listUsers, addUser, removeUser, removeAllUsers };
});
