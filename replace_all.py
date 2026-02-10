
import os

replacements = [
    ("NoFxAiOS/nofx", "xbcvv/nofx-0210"),
    ("ghcr.io/nofxaios/nofx/nofx-backend:stable", "xbcvv/nofx-backend:latest"),
    ("ghcr.io/nofxaios/nofx/nofx-frontend:stable", "xbcvv/nofx-frontend:latest"),
    ("https://raw.githubusercontent.com/NoFxAiOS/nofx/release/stable", "https://raw.githubusercontent.com/xbcvv/nofx-0210/main"),
]

files = [
    r"docker-compose.stable.yml",
    r"web\src\constants\branding.ts",
    r"web\src\components\landing\FooterSection.tsx",
    r"web\src\components\faq\FAQLayout.tsx",
    r"web\src\components\faq\FAQContent.tsx",
    r"scripts\pr-fix.sh",
    r"scripts\pr-check.sh",
    r"install-stable.sh",
]

base_dir = r"d:\qdw\trae\nofx\nofx-0210"

for rel_path in files:
    file_path = os.path.join(base_dir, rel_path)
    if not os.path.exists(file_path):
        print(f"Skipping {file_path} (not found)")
        continue
        
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
            
        new_content = content
        for old, new in replacements:
            new_content = new_content.replace(old, new)
        
        if content != new_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
            print(f"Updated {rel_path}")
        else:
            print(f"No changes in {rel_path}")
            
    except Exception as e:
        print(f"Error processing {rel_path}: {e}")
