<template>
  <v-container>
    <v-row>
      <v-btn :class="{ 'v-btn--active': mode === 'pan' }" @click="setMode('pan')">Pan</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'drag' }" @click="setMode('drag')">Drag</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'add' }" @click="setMode('add')">Add Node</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'remove' }" @click="setMode('remove')">Remove Node/Edge</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'edit' }" @click="setMode('edit')">Edit Node/Edge</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'addEdge' }" @click="setMode('addEdge')">Add Edge</v-btn>
    </v-row>
    <svg
      ref="svg"
      :width="width"
      :height="height"
      :class="{'graph-svg': true, 'pan-mode': mode === 'pan', 'default-mode': mode !== 'pan'}"
      @mousedown="startInteraction"
      @mouseup="stopInteraction"
      @mouseleave="stopInteraction"
      @mousemove="handleMouseOver"
      @click="handleSvgClick"
    >
      <filter id="highlight">
        <feGaussianBlur result="blurOut" in="SourceGraphic" stdDeviation="3"></feGaussianBlur>
        <feFlood flood-color="yellow" flood-opacity="0.5"></feFlood>
        <feComposite in2="blurOut" operator="in"></feComposite>
        <feMerge>
          <feMergeNode></feMergeNode>
          <feMergeNode in="SourceGraphic"></feMergeNode>
        </feMerge>
      </filter>
      <g :transform="`translate(${panX}, ${panY})`">
        <Edge
          v-for="edge in validEdges"
          :key="'edge-' + edge.id"
          :from="findNode(edge.from)"
          :to="findNode(edge.to)"
          :edge="edge"
          :highlighted="(selectedEdgeId === edge.id || highlightedEdgeId === edge.id) && (mode === 'edit' || mode === 'remove')"
          @select="handleEdgeSelect(edge)"
        />
      </g>
      <g :transform="`translate(${panX}, ${panY})`">
        <Node
          v-for="node in nodes"
          :key="'node-' + node.id"
          :id="node.id"
          :x="node.x"
          :y="node.y"
          :selected="node.id === selectedNodeId || node.id === edgeStartNode"
          :highlighted="node.id === hoveredNodeId || node.id === selectedNodeIdInEditMode"
          :mode="mode"
          @select="selectNode"
          @mousedown.stop="startNodeDrag($event, node.id)"
        />
      </g>
      <line
        v-if="mode === 'addEdge' && edgeStartNode && hoveredNodeId"
        :x1="findNode(edgeStartNode).x + panX"
        :y1="findNode(edgeStartNode).y + panY"
        :x2="findNode(hoveredNodeId).x + panX"
        :y2="findNode(hoveredNodeId).y + panY"
        stroke="black"
        opacity="0.5"
        stroke-dasharray="5,5"
      />
      <circle
        v-if="mode === 'add' && addNodeHovered"
        :cx="addNodeHovered.x + panX"
        :cy="addNodeHovered.y + panY"
        r="20"
        fill="transparent"
        stroke="black"
        opacity="0.5"
        stroke-dasharray="5,5"
      />
    </svg>
  </v-container>
</template>

<script>
import { throttle } from 'lodash';
import Node from './Node.vue';
import Edge from './Edge.vue';

export default {
  components: {
    Node,
    Edge,
  },
  data() {
    return {
      nodes: [
        { id: 'S', x: 50, y: 50 },
        { id: 'A', x: 150, y: 50 },
        // ... other nodes
      ],
      edges: [
        { id: 'S-A', from: 'S', to: 'A', weight: 4 },
        { id: 'S-B', from: 'S', to: 'B', weight: 3 },
        // ... other edges
      ],
      selectedNodeId: null,
      hoveredNodeId: null,
      selectedEdgeId: null,
      highlightedEdgeId: null,
      selectedNodeIdInEditMode: null,
      draggingNodeId: null,
      isPanning: false,
      offsetX: 0,
      offsetY: 0,
      panX: 0,
      panY: 0,
      panStartX: 0,
      panStartY: 0,
      width: 800,
      height: 600,
      mode: 'pan',
      creatingEdge: null,
      edgeStartNode: null,
      addNodeHovered: null,
    };
  },
  computed: {
    validEdges() {
      return this.edges.filter(edge => this.findNode(edge.from) && this.findNode(edge.to));
    }
  },
  mounted() {
    this.throttledMouseMove = throttle(this.onMouseMove, 16);
    window.addEventListener('mousemove', this.throttledMouseMove);
    window.addEventListener('mouseup', this.stopInteraction);
  },
  beforeDestroy() {
    window.removeEventListener('mousemove', this.throttledMouseMove);
    window.removeEventListener('mouseup', this.stopInteraction);
  },
  methods: {
    findNode(id) {
      return this.nodes.find(node => node.id === id);
    },
    selectNode(id) {
      if (this.mode === 'remove') {
        this.removeNode(id);
      } else if (this.mode === 'edit') {
        this.selectedNodeIdInEditMode = id;
        const newId = prompt('Enter new ID', id);
        if (newId && !this.findNode(newId)) {
          this.updateNodeId(id, newId);
        }
        this.selectedNodeIdInEditMode = null;
      } else if (this.mode === 'addEdge') {
        if (!this.edgeStartNode) {
          this.edgeStartNode = id;
          this.selectedNodeId = id;
        } else {
          if (this.edgeStartNode !== id) {
            const existingEdge = this.edges.find(edge => 
              (edge.from === this.edgeStartNode && edge.to === id) ||
              (edge.from === id && edge.to === this.edgeStartNode)
            );
            if (!existingEdge) {
              const weight = prompt('Enter weight for new edge');
              if (weight) {
                this.addEdge(this.edgeStartNode, id, parseInt(weight));
                this.edgeStartNode = null;
                this.selectedNodeId = null;
              }
            } else {
              alert('An edge already exists between these nodes.');
            }
          } else {
            alert('Cannot connect a node to itself. Select a different node.');
            this.edgeStartNode = null;
            this.selectedNodeId = null;
          }
        }
      } else if (this.mode === 'pan' || this.mode === 'drag') {
        this.selectedNodeId = id;
      }
    },
    handleEdgeSelect(edge) {
      if (this.mode === 'remove') {
        this.removeEdge(edge.id);
      } else if (this.mode === 'edit') {
        this.selectedEdgeId = edge.id;
        const newWeight = prompt('Enter new weight', edge.weight);
        if (newWeight) {
          this.updateEdgeWeight(edge.id, parseInt(newWeight));
        }
        this.selectedEdgeId = null;
      }
    },
    startNodeDrag(event, nodeId) {
      if (this.mode === 'drag') {
        event.preventDefault();
        this.draggingNodeId = nodeId;
        const node = this.findNode(nodeId);
        if (node) {
          this.offsetX = event.clientX - node.x - this.panX;
          this.offsetY = event.clientY - node.y - this.panY;
        }
      }
    },
    startInteraction(event) {
      if (this.mode === 'pan') {
        this.isPanning = true;
        this.panStartX = event.clientX;
        this.panStartY = event.clientY;
      }
      if (this.mode === 'drag' && this.draggingNodeId) {
        this.selectedNodeId = this.draggingNodeId;
      }
    },
    stopInteraction() {
      this.isPanning = false;
      this.draggingNodeId = null;
    },
    onMouseMove(event) {
      if (this.draggingNodeId !== null && this.mode === 'drag') {
        const node = this.findNode(this.draggingNodeId);
        if (node) {
          node.x = event.clientX - this.offsetX - this.panX;
          node.y = event.clientY - this.offsetY - this.panY;
        }
      }
      if (this.isPanning) {
        const dx = event.clientX - this.panStartX;
        const dy = event.clientY - this.panStartY;
        this.panX += dx;
        this.panY += dy;
        this.panStartX = event.clientX;
        this.panStartY = event.clientY;
      }
    },
    handleMouseOver(event) {
      const svgRect = this.$refs.svg.getBoundingClientRect();
      const x = event.clientX - svgRect.left - this.panX;
      const y = event.clientY - svgRect.top - this.panY;
      const hoveredNode = this.nodes.find(node => Math.abs(node.x - x) < 20 && Math.abs(node.y - y) < 20);

      if (this.mode === 'addEdge' || this.mode === 'edit' || this.mode === 'drag' || this.mode === 'remove') {
        this.hoveredNodeId = hoveredNode ? hoveredNode.id : null;
        if ((this.mode === 'edit' || this.mode === 'remove') && !hoveredNode) {
          const hoveredEdge = this.edges.find(edge => {
            const fromNode = this.findNode(edge.from);
            const toNode = this.findNode(edge.to);
            if (!fromNode || !toNode) return false;
            const dist = Math.abs(
              (toNode.y - fromNode.y) * x - (toNode.x - fromNode.x) * y + toNode.x * fromNode.y - toNode.y * fromNode.x
            ) /
              Math.sqrt(
                Math.pow(toNode.y - fromNode.y, 2) + Math.pow(toNode.x - fromNode.x, 2)
              );
            return dist < 5;
          });
          this.highlightedEdgeId = hoveredEdge ? hoveredEdge.id : null;
        } else {
          this.highlightedEdgeId = null;
        }
      } else if (this.mode === 'add') {
        this.addNodeHovered = { x, y };
      } else {
        this.hoveredNodeId = null;
        this.addNodeHovered = null;
        this.highlightedEdgeId = null;
      }
      if (this.draggingNodeId === null) { // Ensure no conflict with dragging
        this.selectedNodeId = this.hoveredNodeId ? this.hoveredNodeId : null;
      }
    },
    handleSvgClick(event) {
      if (this.mode === 'add') {
        const svgRect = this.$refs.svg.getBoundingClientRect();
        const x = event.clientX - svgRect.left - this.panX;
        const y = event.clientY - svgRect.top - this.panY;
        const newId = prompt('Enter new ID');
        if (newId && !this.findNode(newId)) {
          this.addNode(newId, x, y);
        }
      } else {
        if(this.mode !== 'addEdge') {
          this.selectedNodeId = null;
          this.edgeStartNode = null;
        }
      }
      this.addNodeHovered = null;
    },
    addNode(id, x, y) {
      if (!this.findNode(id)) {
        this.nodes.push({ id, x, y });
      }
    },
    removeNode(id) {
      this.nodes = this.nodes.filter(node => node.id !== id);
      this.edges = this.edges.filter(edge => edge.from !== id && edge.to !== id);
    },
    addEdge(from, to, weight) {
      if (!this.edges.find(edge => (edge.from === from && edge.to === to) || (edge.from === to && edge.to === from))) {
        this.edges.push({ id: `${from}-${to}`, from, to, weight });
      }
    },
    removeEdge(id) {
      this.edges = this.edges.filter(edge => edge.id !== id);
    },
    updateNodeId(oldId, newId) {
      const node = this.findNode(oldId);
      if (node) {
        node.id = newId;
        this.edges.forEach(edge => {
          if (edge.from === oldId) edge.from = newId;
          if (edge.to === oldId) edge.to = newId;
        });
      }
    },
    updateEdgeWeight(id, newWeight) {
      const edge = this.edges.find(edge => edge.id === id);
      if (edge) {
        edge.weight = newWeight;
      }
    },
    setMode(mode) {
      this.mode = mode;
      this.isPanning = false;
      this.draggingNodeId = null;
      this.edgeStartNode = null;
      this.selectedNodeId = null;
      this.hoveredNodeId = null;
      this.highlightedEdgeId = null;
      this.addNodeHovered = null;
    },
  }
};
</script>

<style scoped>
.graph-svg {
  border: 1px solid black;
  transition: transform 0.1s ease, cursor 0.1s ease;
  transform-origin: 0 0;
}
.pan-mode {
  cursor: move;
}
.default-mode {
  cursor: default;
}
.v-btn--active {
  background-color: #1976D2 !important;
  color: white !important;
}
</style>