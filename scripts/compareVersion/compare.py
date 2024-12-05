import subprocess
import re
import os

def get_highest_patch_versions():
    result = subprocess.run(['git', 'tag'], capture_output=True, text=True, check=True)
    tags = result.stdout.strip().split('\n')
    
    tag_pattern = re.compile(r'^v(\d+)\.(\d+)\.(\d+)$')
    
    highest_patch_versions = {}
    
    for tag in tags:
        match = tag_pattern.match(tag)
        if match:
            major, minor, patch = map(int, match.groups())
            if major >= 1:
                major_minor = f"{major}.{minor}"
            # Ensure the value stored is a tuple of (patch, tag)
                if major_minor not in highest_patch_versions or patch > highest_patch_versions[major_minor][0]:
                    highest_patch_versions[major_minor] = (patch, tag)
    
    return highest_patch_versions

def switch_to_tag_and_run_script(tag):
    subprocess.run(['git', 'checkout', tag], check=True)
    os.chdir('..')  # Move back to root where scripts dir is located
    subprocess.run(['go', 'run', 'compare.go',"terraform-"+tag], check=True)
    os.chdir('terraform')  # Move back to terraform dir


def switch_to_tag_and_run_script2(tag):
    subprocess.run(['git', 'checkout', tag], check=True)
    os.chdir('..')  # Move back to root where scripts dir is located
    subprocess.run(['go', 'run', 'compare.go',"opentofu-"+tag], check=True)
    os.chdir('opentofu')  # Move back to terraform dir

def main():
    subprocess.run(['rm', 'versions.json'], check=True)
    terraform_dir = 'terraform'
    print("terraform 兼容情况:")
    if not os.path.exists(terraform_dir):
        subprocess.run(['git', 'clone', 'https://github.com/hashicorp/terraform.git', terraform_dir])
    
    os.chdir(terraform_dir)
    
    highest_patch_versions = get_highest_patch_versions()
    
    for major_minor, (patch, tag) in highest_patch_versions.items():
        print(f"v{major_minor}.{patch}")
        switch_to_tag_and_run_script(tag)
    
    os.chdir('..')

    terraform_dir = 'opentofu'
    print("opentofu 兼容情况:")
    if not os.path.exists(terraform_dir):
        subprocess.run(['git', 'clone', 'https://github.com/opentofu/opentofu.git', terraform_dir])
    
    os.chdir(terraform_dir)
    
    highest_patch_versions = get_highest_patch_versions()
    
    for major_minor, (patch, tag) in highest_patch_versions.items():
        print(f"v{major_minor}.{patch}")
        switch_to_tag_and_run_script2(tag)
    
    os.chdir('..')
    os.chdir('analys')
    subprocess.run(['go', 'run', 'analys.go',"opentofu-"+tag], check=True)


if __name__ == "__main__":
    main()
