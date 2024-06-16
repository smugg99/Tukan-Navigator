<template>
  <v-container>
    <v-row>
      <v-btn :class="{ 'v-btn--active': mode === 'pan' }" @click="setMode('pan')">Pan</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'drag' }" @click="setMode('drag')">Drag</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'add' }" @click="setMode('add')">Add Node</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'remove' }" @click="setMode('remove')">Remove Node/Edge</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'edit' }" @click="setMode('edit')">Edit Node/Edge</v-btn>
      <v-btn :class="{ 'v-btn--active': mode === 'addEdge' }" @click="setMode('addEdge')">Add Edge</v-btn>
      <v-btn @click="getGraphRelationsAndPositions">Print Graph Relations and Positions</v-btn>
      <v-btn @click="getGraphRelations">Print Graph Relations</v-btn>
      <v-btn @click="startAnimation">Start Animation</v-btn>
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
      <circle
        v-if="isPathSet"
        :cx="currentPosition.x + panX"
        :cy="currentPosition.y + panY"
        r="10"
        fill="blue"
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
      ],
      edges: [
        { id: 'S-A', from: 'S', to: 'A', weight: 4 },
        { id: 'S-B', from: 'S', to: 'B', weight: 3 },
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
      isPathSet: false,
      animatedPath: [],
      currentPosition: { x: 0, y: 0 },
      pathIndex: 0,
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
        if (newId && newId !== id) {
          if (!this.findNode(newId)) {
            this.updateNodeId(id, newId);
          } else {
            alert('Duplicate ID detected. Please enter a unique ID.');
          }
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
        this.panStartX = event.clientX - this.panX;
        this.panStartY = event.clientY - this.panY;
      }
    },
    stopInteraction() {
      this.isPanning = false;
      this.draggingNodeId = null;
    },
    onMouseMove(event) {
      if (this.isPanning) {
        this.panX = event.clientX - this.panStartX;
        this.panY = event.clientY - this.panStartY;
      } else if (this.draggingNodeId) {
        const node = this.findNode(this.draggingNodeId);
        if (node) {
          node.x = event.clientX - this.offsetX - this.panX;
          node.y = event.clientY - this.offsetY - this.panY;
        }
      }
    },
    handleMouseOver(event) {
      if (this.mode === 'add') {
        this.addNodeHovered = { x: event.clientX - this.panX, y: event.clientY - this.panY };
      } else {
        this.hoveredNodeId = this.nodes.find(node => 
          Math.abs(event.clientX - node.x - this.panX) < 10 &&
          Math.abs(event.clientY - node.y - this.panY) < 10
        )?.id || null;
      }
    },
    handleSvgClick(event) {
      if (this.mode === 'add' && this.addNodeHovered) {
        const id = prompt('Enter node ID');
        if (id) {
          this.addNode(id, this.addNodeHovered.x, this.addNodeHovered.y);
        }
        this.addNodeHovered = null;
      }
    },
    addNode(id, x, y) {
      if (this.findNode(id)) {
        alert('A node with this ID already exists.');
      } else {
        this.nodes.push({ id, x, y });
      }
    },
    removeNode(id) {
      this.nodes = this.nodes.filter(node => node.id !== id);
      this.edges = this.edges.filter(edge => edge.from !== id && edge.to !== id);
    },
    addEdge(from, to, weight) {
      const id = `${from}-${to}`;
      if (!this.edges.some(edge => edge.id === id)) {
        this.edges.push({ id, from, to, weight });
      } else {
        alert('This edge already exists.');
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
    updateEdgeWeight(id, weight) {
      const edge = this.edges.find(edge => edge.id === id);
      if (edge) edge.weight = weight;
    },
    setMode(newMode) {
      this.mode = newMode;
    },
    getGraphRelationsAndPositions() {
      const relationsAndPositions = {
        nodes: this.nodes.map(node => ({ id: node.id, x: node.x, y: node.y })),
        edges: this.edges.map(edge => ({ from: edge.from, to: edge.to, weight: edge.weight }))
      };
      console.log(JSON.stringify(relationsAndPositions, null, 2));
    },
    getGraphRelations() {
      const relations = this.edges.map(edge => ({
        from: edge.from,
        to: edge.to,
        weight: edge.weight
      }));
      console.log(JSON.stringify(relations, null, 2));
    },
    async startAnimation() {
      try {
        const response = await fetch('http://localhost:8000/api/v1/graphs/shortest-path', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
              nodes: this.nodes.map(node => node.id),
              edges: this.edges.map(edge => ({
                from: edge.from,
                to: edge.to,
                weight: edge.weight
              }))
            },
          )
        });
        console.log(response);
        const data = await response.json();
        if (data.path) {
          this.animatedPath = data.path.map(nodeId => {
            const node = this.findNode(nodeId);
            return { x: node.x, y: node.y };
          });
          this.isPathSet = true;
          this.pathIndex = 0;
          this.animateTukan();
        } else {
          console.error('Path data is not valid:', data);
        }
      } catch (error) {
        console.error('Error fetching path:', error);
      }
    },
    animateTukan() {
      if (this.pathIndex < this.animatedPath.length) {
        this.currentPosition = this.animatedPath[this.pathIndex];
        this.pathIndex++;
        setTimeout(this.animateTukan, 1000); // Adjust the delay as needed
      }
    },
  }
};
</script>

<style scoped>
.graph-svg {
  border: 1px solid #ccc;
  background-color: #f9f9f9;
  cursor: grab;
}

.graph-svg.pan-mode {
  cursor: grab;
}

.graph-svg.default-mode {
  cursor: default;
}
</style>