#!/usr/bin/env node

/**
 * Generate PDF Resume from HTML Template
 *
 * This script uses Puppeteer to convert the resume HTML template
 * into a professional 2-page PDF document.
 *
 * Usage: node scripts/generate-resume-pdf.js
 */

const puppeteer = require('puppeteer');
const path = require('path');
const fs = require('fs');

async function generateResumePDF() {
  console.log('🚀 Starting PDF generation...');

  // Paths
  const htmlPath = path.join(__dirname, '../web/templates/pages/resume.html');
  const outputPath = path.join(__dirname, '../web/static/resume.pdf');

  // Check if HTML file exists
  if (!fs.existsSync(htmlPath)) {
    console.error(`❌ Error: Resume template not found at ${htmlPath}`);
    process.exit(1);
  }

  // Read HTML content
  const htmlContent = fs.readFileSync(htmlPath, 'utf-8');

  let browser;
  try {
    // Launch browser
    console.log('📖 Launching browser...');
    browser = await puppeteer.launch({
      headless: 'new',
      args: ['--no-sandbox', '--disable-setuid-sandbox']
    });

    const page = await browser.newPage();

    // Set content
    console.log('📄 Loading resume content...');
    await page.setContent(htmlContent, {
      waitUntil: 'networkidle0'
    });

    // Generate PDF
    console.log('🖨️  Generating PDF...');
    await page.pdf({
      path: outputPath,
      format: 'Letter',
      printBackground: true,
      margin: {
        top: '0.5in',
        right: '0.5in',
        bottom: '0.5in',
        left: '0.5in'
      },
      preferCSSPageSize: true
    });

    console.log(`✅ PDF generated successfully: ${outputPath}`);
    console.log(`📊 File size: ${(fs.statSync(outputPath).size / 1024).toFixed(2)} KB`);

  } catch (error) {
    console.error('❌ Error generating PDF:', error);
    process.exit(1);
  } finally {
    if (browser) {
      await browser.close();
    }
  }
}

// Run the generator
generateResumePDF().catch(error => {
  console.error('❌ Fatal error:', error);
  process.exit(1);
});
