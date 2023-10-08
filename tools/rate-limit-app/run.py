#!/usr/bin/env python
import time
import typer
import requests
import random
import socket
import struct
from rich import print
from rich.progress import Progress, SpinnerColumn, TimeElapsedColumn
from rich.columns import Columns
from rich.console import Console

def main(
        url: str = typer.Argument(..., help='URL to test'),
        method: str = typer.Option('GET', help='HTTP method to use'),
        delay: float = typer.Option(1, help='Delay in seconds between requests'),
        limit: int = typer.Option(100, help='Number of requests to make (0: unlimited)'),
        xforwardedfor_random: bool = typer.Option(False, help='Should a random X-Forwarded-For header be injected')
    ):
    print(f'URL:                     {url}')
    print(f'Method:                  {method}')
    print(f'Delay:                   {delay} s')
    print(f'Limit:                   {limit}')
    print(f'Random X-Forwarded-For:  {xforwardedfor_random}')
    i = 0
    console = Console()
    session = requests.Session()
    headers = {}
    with Progress(SpinnerColumn(), *Progress.get_default_columns(), TimeElapsedColumn(), console=console) as progress:
        task = progress.add_task("[yellow]Sending requests...", total=(limit if limit > 0 else None))
        while not progress.finished:
            if xforwardedfor_random:
                headers["X-Forwarded-For"] = str(socket.inet_ntoa(struct.pack('>I', random.randint(1, 0xffffffff))))
            request = requests.Request(method=method, url=url, headers=headers).prepare()
            resp = session.send(request)
            progress.console.print(Columns([f"[yellow]#{i}", f"[blue]{int(round(resp.elapsed.total_seconds()*1000, 0))}ms", f"[{'green' if resp.status_code == 200 else 'red'}]{resp.status_code}"]))
            progress.advance(task)
            time.sleep(delay)
            i += 1
if __name__ == '__main__':
    typer.run(main)