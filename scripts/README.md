# Resume PDF Generator

This directory contains scripts for generating a PDF version of Maria Lucena's resume from the HTML template.

## Files

- `generate-resume-pdf.js` - Node.js script that uses Puppeteer to convert HTML to PDF

## Usage

### Generate Resume PDF

```bash
# Using npm script (recommended)
npm run generate:resume

# Using Makefile
make resume

# Direct execution
node scripts/generate-resume-pdf.js
```

## Output

The generated PDF is saved to:
```
web/static/resume.pdf
```

This file is automatically referenced by the "Download Resume" button on the homepage.

## Resume Template

The HTML resume template is located at:
```
web/templates/pages/resume.html
```

To update the resume content, edit this HTML file and regenerate the PDF.

## Technical Details

- **PDF Size**: ~270 KB
- **Format**: US Letter (8.5" x 11")
- **Margins**: 0.5 inches on all sides
- **Layout**: Optimized for 2 pages
- **Print Settings**: Background colors and graphics enabled

## Dependencies

- `puppeteer` - Headless Chrome for PDF generation
- Automatically installed via `npm install`

## Customization

To customize the PDF output, modify the `page.pdf()` options in `generate-resume-pdf.js`:

```javascript
await page.pdf({
  path: outputPath,
  format: 'Letter',      // Paper format
  printBackground: true, // Include colors/backgrounds
  margin: {              // Page margins
    top: '0.5in',
    right: '0.5in',
    bottom: '0.5in',
    left: '0.5in'
  }
});
```

## Resume Content Sections

The generated PDF includes:

1. **Header** - Name, title, contact information
2. **Professional Summary** - Overview of experience and expertise
3. **Core Competencies** - Key skills organized by category
4. **Professional Experience** - Detailed work history with achievements
5. **Key Projects** - Notable projects with technologies
6. **Patents** - List of granted and pending patents
7. **Speaking Engagements** - Conference talks and presentations
8. **Education** - Academic credentials

## Maintenance

When updating the portfolio website with new information, remember to:

1. Update the resume HTML template (`web/templates/pages/resume.html`)
2. Regenerate the PDF (`make resume`)
3. Commit both files to version control
