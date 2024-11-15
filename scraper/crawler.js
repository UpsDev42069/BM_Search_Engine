import fs from 'fs';
import { JSDOM } from 'jsdom';
import { URL } from 'url';
import { indexWebPage } from './indexer.js';
// If Node.js < v18, uncomment the following line after installing node-fetch
// import fetch from 'node-fetch';

const visitedUrls = new Set();

const writeStream = fs.createWriteStream('crawledData.json');
writeStream.write('[\n');
let isFirstEntry = true;


function delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function crawlPage(pageUrl) {

    const cleanUrl = new URL(pageUrl);
    cleanUrl.hash = '';

    if (visitedUrls.has(cleanUrl.href)) {
        return;
    }

    visitedUrls.add(cleanUrl.href);

    try {
        const response = await fetch(cleanUrl.href);
        if (!response.ok) {
            console.error(`Failed to fetch ${cleanUrl.href}: ${response.statusText}`);
            return;
        }

        const html = await response.text();
        const { window } = new JSDOM(html);
        const document = window.document;

        const title = document.querySelector('.mw-first-heading').textContent.trim();


        const mainContent = Array.from(document.querySelectorAll('.mw-parser-output p, .mw-parser-output h1, .mw-parser-output h2, .mw-parser-output h3'))
            .map(element => element.textContent.trim())
            .join(' ')
            .replace(/[\t\n\r]+/g, ' ')
            .trim();

        console.log(`Crawling and indexing: ${cleanUrl.href}`);


        const pageData = {url: cleanUrl.href, title: title, content: mainContent};

        if (!isFirstEntry){
            writeStream.write(',\n');
        
        } else {
            isFirstEntry = false;
        }
        writeStream.write(JSON.stringify(pageData));

        const links = Array.from(document.querySelectorAll('a')).map(link => link.getAttribute('href'));
        for (let href of links) {
            if (href && href.startsWith('/wiki/') && !href.includes(':')) {
                const resolvedUrl = new URL(href, cleanUrl.origin).href;
                const normalizedUrl = resolvedUrl.split('#')[0];
                
                if (!visitedUrls.has(normalizedUrl)) {
                    await delay(1000); // Delay between requests
                    await crawlPage(normalizedUrl);
                }
            }

        }
    } catch (error) {
        console.error(`Error crawling ${cleanUrl.href}:`, error);
    }
}

// Invoke crawlPage once and handle the results
crawlPage('https://en.wikipedia.org/wiki/Search_engine')