import knexConfig from "./knexfile.js";
import knex from "knex";
import fs from "fs/promises";

async function importData() {
  const knexPostgres = knex(knexConfig.postgres);
  try {
    const data = await fs.readFile("crawledData.json");
    const crawledData = JSON.parse(data);

    if (!Array.isArray(crawledData)) {
      throw new Error("Invalid data format");
    }

    const formattedData = crawledData.map((page) => ({
      url: page.url,
      title: page.title,
      content: page.content,
    }));

    await knexPostgres("pages").insert(formattedData).onConflict("title", "url").ignore();

    console.log("Migration successful");
  } catch (error) {
    console.error("Migration failed: ", error);
  } finally {
    await knexPostgres.destroy();
  }
}

importData();
