#!/usr/bin/env python

from setuptools import find_packages, setup

setup(
    name='whichip',
    version='0.1.1',
    description='discover (IoT) device\'s IP(s) in local network',
    long_description=(
        open('README.md').read()
    ),
    long_description_content_type="text/markdown",
    author='observerss',
    url='https://github.com/observerss/whichip',
    py_modules=['whichip'],
    packages=find_packages(where="."),
    entry_points={
        'console_scripts': ['whichip=whichip:main'],
    },
    license='MIT',
    classifiers=[
        'Development Status :: 5 - Production/Stable',
        'Intended Audience :: Developers',
        'Natural Language :: English',
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python',
        'Programming Language :: Python :: 3 :: Only',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
        'Programming Language :: Python :: 3.9',
        'Programming Language :: Python :: 3.10',
    ],
    python_requires='>=3.6'
)
