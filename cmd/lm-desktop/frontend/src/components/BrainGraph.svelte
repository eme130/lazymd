<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import * as d3 from 'd3';
  import { forceSimulation, forceLink, forceManyBody, forceCenter } from 'd3-force';
  import { GetGraph, OpenFile } from '../../wailsjs/go/wailsplugin/App';
  import { onGraphUpdated } from '../lib/events';

  interface GraphNode extends d3.SimulationNodeDatum {
    name: string;
    path: string;
    linkCount: number;
  }

  interface GraphEdge {
    source: string | GraphNode;
    target: string | GraphNode;
  }

  let canvas: HTMLCanvasElement;
  let simulation: d3.Simulation<GraphNode, GraphEdge>;
  let transform = d3.zoomIdentity;
  let nodes: GraphNode[] = [];
  let links: GraphEdge[] = [];

  async function loadGraph() {
    const data = await GetGraph();
    if (!data?.nodes) return;
    nodes = data.nodes.map((n: any) => ({ ...n }));
    links = data.edges?.map((e: any) => ({ ...e })) || [];
    setupSimulation();
  }

  function setupSimulation() {
    if (!canvas) return;
    const width = canvas.width;
    const height = canvas.height;
    simulation?.stop();
    simulation = forceSimulation<GraphNode>(nodes)
      .force('link', forceLink<GraphNode, GraphEdge>(links).id((d) => d.name).distance(80))
      .force('charge', forceManyBody().strength(-200))
      .force('center', forceCenter(width / 2, height / 2))
      .on('tick', draw);
  }

  function draw() {
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;
    const width = canvas.width;
    const height = canvas.height;
    ctx.save();
    ctx.clearRect(0, 0, width, height);
    ctx.translate(transform.x, transform.y);
    ctx.scale(transform.k, transform.k);

    ctx.strokeStyle = '#3b4261';
    ctx.lineWidth = 1;
    for (const link of links) {
      const source = link.source as GraphNode;
      const target = link.target as GraphNode;
      if (source.x == null || target.x == null) continue;
      ctx.beginPath();
      ctx.moveTo(source.x, source.y!);
      ctx.lineTo(target.x, target.y!);
      ctx.stroke();
    }

    for (const node of nodes) {
      if (node.x == null) continue;
      const radius = 4 + Math.min(node.linkCount * 2, 16);
      ctx.beginPath();
      ctx.arc(node.x, node.y!, radius, 0, 2 * Math.PI);
      ctx.fillStyle = '#7aa2f7';
      ctx.fill();
      ctx.fillStyle = '#c0caf5';
      ctx.font = '11px sans-serif';
      ctx.textAlign = 'center';
      ctx.fillText(node.name, node.x, node.y! + radius + 14);
    }
    ctx.restore();
  }

  function handleClick(event: MouseEvent) {
    const rect = canvas.getBoundingClientRect();
    const x = (event.clientX - rect.left - transform.x) / transform.k;
    const y = (event.clientY - rect.top - transform.y) / transform.k;
    for (const node of nodes) {
      if (node.x == null) continue;
      const radius = 4 + Math.min(node.linkCount * 2, 16);
      const dx = x - node.x;
      const dy = y - node.y!;
      if (dx * dx + dy * dy < radius * radius) {
        OpenFile(node.path);
        return;
      }
    }
  }

  onMount(() => {
    const resizeObserver = new ResizeObserver(() => {
      if (canvas) {
        canvas.width = canvas.parentElement?.clientWidth || 400;
        canvas.height = canvas.parentElement?.clientHeight || 400;
        draw();
      }
    });
    resizeObserver.observe(canvas.parentElement!);
    d3.select(canvas).call(
      d3.zoom<HTMLCanvasElement, unknown>()
        .scaleExtent([0.2, 5])
        .on('zoom', (event) => { transform = event.transform; draw(); })
    );
    loadGraph();
    onGraphUpdated(() => loadGraph());
    return () => { simulation?.stop(); resizeObserver.disconnect(); };
  });

  onDestroy(() => { simulation?.stop(); });
</script>

<div class="brain-container">
  <canvas bind:this={canvas} on:click={handleClick}></canvas>
</div>

<style>
  .brain-container { width: 100%; height: 100%; overflow: hidden; }
  canvas { display: block; cursor: grab; }
  canvas:active { cursor: grabbing; }
</style>
