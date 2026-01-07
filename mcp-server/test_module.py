#!/usr/bin/env python3
"""
WeKnora MCP Server Module Test Script

Test various startup methods and functionality of the module
"""

import os
import subprocess
import sys
from pathlib import Path


def test_imports():
    """Test module imports"""
    print("=== Testing Module Imports ===")

    try:
        # Test basic dependencies
        import mcp

        print("✓ mcp module imported successfully")

        import requests

        print("✓ requests module imported successfully")

        # Test main module
        import weknora_mcp_server

        print("✓ weknora_mcp_server module imported successfully")

        # Test package imports
        from weknora_mcp_server import WeKnoraClient, run

        print("✓ WeKnoraClient and run function imported successfully")

        # Test main entry point
        import main

        print("✓ main module imported successfully")

        return True

    except ImportError as e:
        print(f"✗ Import failed: {e}")
        return False


def test_environment():
    """Test environment configuration"""
    print("\n=== Testing Environment Configuration ===")

    base_url = os.getenv("WEKNORA_BASE_URL")
    api_key = os.getenv("WEKNORA_API_KEY")

    print(f"WEKNORA_BASE_URL: {base_url or 'Not set (will use default value)'}")
    print(f"WEKNORA_API_KEY: {'Set' if api_key else 'Not set'}")

    if not base_url:
        print("Tip: You can set the WEKNORA_BASE_URL environment variable")

    if not api_key:
        print("Tip: It is recommended to set the WEKNORA_API_KEY environment variable")

    return True


def test_client_creation():
    """Test client creation"""
    print("\n=== Testing Client Creation ===")

    try:
        from weknora_mcp_server import WeKnoraClient

        base_url = os.getenv("WEKNORA_BASE_URL", "http://localhost:8080/api/v1")
        api_key = os.getenv("WEKNORA_API_KEY", "test_key")

        client = WeKnoraClient(base_url, api_key)
        print("✓ WeKnoraClient created successfully")

        # Check client attributes
        assert client.base_url == base_url
        assert client.api_key == api_key
        print("✓ Client configuration is correct")

        return True

    except Exception as e:
        print(f"✗ Client creation failed: {e}")
        return False


def test_file_structure():
    """Test file structure"""
    print("\n=== Testing File Structure ===")

    required_files = [
        "__init__.py",
        "main.py",
        "run_server.py",
        "weknora_mcp_server.py",
        "requirements.txt",
        "setup.py",
        "pyproject.toml",
        "README.md",
        "INSTALL.md",
        "LICENSE",
        "MANIFEST.in",
    ]

    missing_files = []
    for file in required_files:
        if Path(file).exists():
            print(f"✓ {file}")
        else:
            print(f"✗ {file} (missing)")
            missing_files.append(file)

    if missing_files:
        print(f"Missing files: {missing_files}")
        return False

    print("✓ All required files exist")
    return True


def test_entry_points():
    """Test entry points"""
    print("\n=== Testing Entry Points ===")

    # Test main.py help option
    try:
        result = subprocess.run(
            [sys.executable, "main.py", "--help"],
            capture_output=True,
            text=True,
            timeout=10,
        )
        if result.returncode == 0:
            print("✓ main.py --help works correctly")
        else:
            print(f"✗ main.py --help failed: {result.stderr}")
            return False
    except subprocess.TimeoutExpired:
        print("✗ main.py --help timed out")
        return False
    except Exception as e:
        print(f"✗ main.py --help error: {e}")
        return False

    # Test environment check
    try:
        result = subprocess.run(
            [sys.executable, "main.py", "--check-only"],
            capture_output=True,
            text=True,
            timeout=10,
        )
        if result.returncode == 0:
            print("✓ main.py --check-only works correctly")
        else:
            print(f"✗ main.py --check-only failed: {result.stderr}")
            return False
    except subprocess.TimeoutExpired:
        print("✗ main.py --check-only timed out")
        return False
    except Exception as e:
        print(f"✗ main.py --check-only error: {e}")
        return False

    return True


def test_package_installation():
    """Test package installation (development mode)"""
    print("\n=== Testing Package Installation ===")

    try:
        # Check if package can be installed in development mode
        result = subprocess.run(
            [sys.executable, "setup.py", "check"],
            capture_output=True,
            text=True,
            timeout=30,
        )

        if result.returncode == 0:
            print("✓ setup.py check passed")
        else:
            print(f"✗ setup.py check failed: {result.stderr}")
            return False

    except subprocess.TimeoutExpired:
        print("✗ setup.py check timed out")
        return False
    except Exception as e:
        print(f"✗ setup.py check error: {e}")
        return False

    return True


def main():
    """Run all tests"""
    print("WeKnora MCP Server Module Test")
    print("=" * 50)

    tests = [
        ("Module Imports", test_imports),
        ("Environment Configuration", test_environment),
        ("Client Creation", test_client_creation),
        ("File Structure", test_file_structure),
        ("Entry Points", test_entry_points),
        ("Package Installation", test_package_installation),
    ]

    passed = 0
    total = len(tests)

    for test_name, test_func in tests:
        try:
            if test_func():
                passed += 1
            else:
                print(f"Test failed: {test_name}")
        except Exception as e:
            print(f"Test exception: {test_name} - {e}")

    print("\n" + "=" * 50)
    print(f"Test Results: {passed}/{total} passed")

    if passed == total:
        print("✓ All tests passed! Module can be used normally.")
        return True
    else:
        print("✗ Some tests failed, please check the errors above.")
        return False


if __name__ == "__main__":
    success = main()
    sys.exit(0 if success else 1)
