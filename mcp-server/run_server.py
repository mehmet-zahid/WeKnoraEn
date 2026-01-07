#!/usr/bin/env python3
"""
WeKnora MCP Server Startup Script
"""

import asyncio
import os
import sys


def check_environment():
    """Check environment configuration"""
    base_url = os.getenv("WEKNORA_BASE_URL")
    api_key = os.getenv("WEKNORA_API_KEY")

    if not base_url:
        print(
            "Warning: WEKNORA_BASE_URL environment variable not set, using default: http://localhost:8080/api/v1"
        )

    if not api_key:
        print("Warning: WEKNORA_API_KEY environment variable not set")

    print(f"WeKnora Base URL: {base_url or 'http://localhost:8080/api/v1'}")
    print(f"API Key: {'Set' if api_key else 'Not set'}")


def main():
    """Main function"""
    print("Starting WeKnora MCP Server...")
    check_environment()

    try:
        from weknora_mcp_server import run

        asyncio.run(run())
    except ImportError as e:
        print(f"Import error: {e}")
        print("Please ensure all dependencies are installed: pip install -r requirements.txt")
        sys.exit(1)
    except KeyboardInterrupt:
        print("\nServer stopped")
    except Exception as e:
        print(f"Server runtime error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
