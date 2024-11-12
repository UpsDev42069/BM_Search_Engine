import { session } from "../stores/session";

export const API_URL = import.meta.env.VITE_API_URL;

export async function login(username: string, password: string) {
  try {
    const response = await fetch(`${API_URL}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
      credentials: "include",
    });
    if (response.ok) {
      await session.checkSession();
    }
    return response;
  } catch (error) {
    console.error("Error logging in", error);
  }
}

export async function logout() {
  try {
    const response = await fetch(`${API_URL}/logout`, {
      credentials: "include",
    });
    if (response.ok) {
      await session.checkSession();
    }
    return response;
  } catch (error) {
    console.error("Error logging out", error);
  }
}

export async function register(username: string, password: string, email: string) {
  try {
    const response = await fetch(`${API_URL}/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password, email }),
      credentials: "include",
    });
    if (response.ok) {
      await session.checkSession();
    }
    return response;
  } catch (error) {
    console.error("Error registering", error);
  }
}

export async function getWeather() {
  try {
    const response = await fetch(`${API_URL}/weather`);
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error getting weather", error);
  }
}

export async function search(query: string) {
  try {
    const res = await fetch(`${API_URL}/search?q=${query}`);
    const data = await res.json();
    return data;
  } catch (error) {
    console.log(error);
  }
}

export async function resetPassword(username: string, oldPassword: string, newPassword: string) {
  try {
    const res = await fetch(`${API_URL}/reset-password`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, oldPassword, newPassword }),
      credentials: "include",
    });
    if (res.ok) {
      return res;
    }
  } catch (error) {
    console.log(error);
  }

}
