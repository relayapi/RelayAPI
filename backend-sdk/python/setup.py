from setuptools import setup, find_packages

setup(
    name="relayapi-sdk",
    version="1.0.0",
    description="RelayAPI Python SDK - A client library for RelayAPI Server",
    author="RelayAPI Team",
    author_email="support@relayapi.com",
    packages=find_packages(),
    install_requires=[
        "requests>=2.25.0",
        "pycryptodome>=3.10.0",
    ],
    python_requires=">=3.7",
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
)