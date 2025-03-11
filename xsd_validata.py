# encoding: utf-8
'''
Created on 2025年2月24日

@author: jingyu.wy
'''
from pathlib import Path
import sys
from lxml import etree


def validate(xsd_content:str, xml_content:str) -> tuple[bool, str]:
    xml_schema = etree.XMLSchema(etree.fromstring(xsd_content))  # @UndefinedVariable
    xml_content = etree.fromstring(xml_content)
    try:
        xml_schema.assertValid(xml_content)
        return True, None
    except etree.DocumentInvalid as e:
        return False, xml_schema.error_log
    
def main(xsd_file:Path, xml_dirs:Path) -> dict:
    results = {}
    with open(xsd_file, 'r') as file:
        xsd_content = file.read()
    for xml_file in xml_dirs.iterdir():
        with open(xml_file, 'r') as file:
            xml_content = file.read()
            xml_content = xml_content.split('\n', 1)[-1]
        xml = etree.fromstring(xml_content)
        api_name = f"{xml.attrib['product']}_{xml.attrib['version']}_{xml.attrib['name']}"
        r = validate(xsd_content, xml_content)
        results[api_name] = {
            "pass_or_not": r[0],
            "error_msgs": "" if r[0] else str(r[1].last_error),
        }
    return results


if __name__ == '__main__':
    xsd_file = Path(sys.argv[1])
    xml_dirs = Path(sys.argv[2])
    results = main(xsd_file, xml_dirs)
    errors = set()
    for k, v in results.items():
        print(f"{k} {v['pass_or_not']} : {v['error_msgs']}")
        if v["error_msgs"]:
            errors.add(v["error_msgs"].split(" ", 1)[-1])

    with open("/tmp/xsd.txt", "w") as file:
        for e in errors:
            file.write(e + "\n")
